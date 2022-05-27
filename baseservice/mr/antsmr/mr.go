package antsmr

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/chuan-fu/Common/baseservice/ants"
	"github.com/chuan-fu/Common/baseservice/mr"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	// ErrCancelWithNil is an error that mapreduce was cancelled with nil.
	ErrCancelWithNil = errors.New("mapreduce cancelled with nil")
	// ErrReduceNoOutput is an error that reduce did not output a value.
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	// ForEachFunc is used to do element processing, but no output.
	ForEachFunc func(item interface{})
	// GenerateFunc is used to let callers send elements into source.
	GenerateFunc func(source chan<- interface{})
	// MapFunc is used to do element processing and write the output to writer.
	MapFunc func(item interface{}, writer Writer)
	// MapperFunc is used to do element processing and write the output to writer,
	// use cancel func to cancel the processing.
	MapperFunc func(item interface{}, writer Writer, cancel func(error))
	// ReducerFunc is used to reduce all the mapping output and write to writer,
	// use cancel func to cancel the processing.
	ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(error))
	// VoidReducerFunc is used to reduce all the mapping output, but no output.
	// Use cancel func to cancel the processing.
	VoidReducerFunc func(pipe <-chan interface{}, cancel func(error))
	// Option defines the method to customize the mapreduce.
	Option func(opts *mapReduceOptions)

	mapperContext struct {
		ctx       context.Context
		mapper    MapFunc
		source    <-chan interface{}
		panicChan *onceChan
		collector chan<- interface{}
		doneChan  <-chan mr.PlaceholderType
		workers   int
	}

	mapReduceOptions struct {
		ctx     context.Context
		workers int
	}

	// Writer interface wraps Write method.
	Writer interface {
		Write(v interface{})
	}
)

// workers为工作队列数量，应<=len(fns)
func FinishWithWorkers(ctx context.Context, workers int, fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}
	if workers > len(fns) {
		workers = len(fns)
	}

	return MapReduceVoid(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}, writer Writer, cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan interface{}, cancel func(error)) {
	}, WithWorkers(workers), WithContext(ctx))
}

// Finish runs fns parallelly, cancelled on any error.
func FinishWithCtx(ctx context.Context, fns ...func() error) error {
	return FinishWithWorkers(ctx, len(fns), fns...)
}

// Finish runs fns parallelly, cancelled on any error.
func Finish(fns ...func() error) error {
	return FinishWithWorkers(nil, len(fns), fns...)
}

func FinishVoidWithWorkers(ctx context.Context, workers int, fns ...func()) {
	if len(fns) == 0 {
		return
	}
	if workers > len(fns) {
		workers = len(fns)
	}

	ForEach(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}) {
		fn := item.(func())
		fn()
	}, WithWorkers(workers), WithContext(ctx))
}

func FinishVoidWithCtx(ctx context.Context, fns ...func()) {
	FinishVoidWithWorkers(ctx, len(fns), fns...)
}

// FinishVoid runs fns parallelly.
func FinishVoid(fns ...func()) {
	FinishVoidWithWorkers(nil, len(fns), fns...)
}

// ForEach maps all elements from given generate but no output.
func ForEach(generate GenerateFunc, mapper ForEachFunc, opts ...Option) {
	options := buildOptions(opts...)
	panicChan := &onceChan{channel: make(chan interface{})}
	source := buildSource(generate, panicChan)
	collector := make(chan interface{}, options.workers)
	done := make(chan mr.PlaceholderType)

	ants.GoGo(func() {
		executeMappers(mapperContext{
			ctx: options.ctx,
			mapper: func(item interface{}, writer Writer) {
				mapper(item)
			},
			source:    source,
			panicChan: panicChan,
			collector: collector,
			doneChan:  done,
			workers:   options.workers,
		})
	})

	for {
		select {
		case v := <-panicChan.channel:
			panic(v)
		case _, ok := <-collector:
			if !ok {
				return
			}
		}
	}
}

// MapReduce maps all elements generated from given generate func,
// and reduces the output elements with given reducer.
func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {
	panicChan := &onceChan{channel: make(chan interface{})}
	source := buildSource(generate, panicChan)
	return mapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
}

// MapReduceChan maps all elements from source, and reduce the output elements with given reducer.
func MapReduceChan(source <-chan interface{}, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {
	panicChan := &onceChan{channel: make(chan interface{})}
	return mapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
}

// MapReduceChan maps all elements from source, and reduce the output elements with given reducer.
func mapReduceWithPanicChan(source <-chan interface{}, panicChan *onceChan, mapper MapperFunc,
	reducer ReducerFunc, opts ...Option) (interface{}, error) {
	options := buildOptions(opts...)
	// output is used to write the final result
	output := make(chan interface{})
	defer func() {
		// reducer can only write once, if more, panic
		for range output {
			panic("more than one element written in reducer")
		}
	}()

	// collector is used to collect data from mapper, and consume in reducer
	collector := make(chan interface{}, options.workers)
	// if done is closed, all mappers and reducer should stop processing
	done := make(chan mr.PlaceholderType)
	writer := newGuardedWriter(options.ctx, output, done)
	var closeOnce sync.Once
	// use atomic.Value to avoid data race
	var retErr mr.AtomicError
	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	ants.GoGo(func() {
		defer func() {
			drain(collector)
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			finish()
		}()

		reducer(collector, writer, cancel)
	})

	ants.GoGo(func() {
		executeMappers(mapperContext{
			ctx: options.ctx,
			mapper: func(item interface{}, w Writer) {
				mapper(item, w, cancel)
			},
			source:    source,
			panicChan: panicChan,
			collector: collector,
			doneChan:  done,
			workers:   options.workers,
		})
	})

	select {
	case <-options.ctx.Done():
		cancel(context.DeadlineExceeded)
		return nil, context.DeadlineExceeded
	case v := <-panicChan.channel:
		panic(v)
	case v, ok := <-output:
		if err := retErr.Load(); err != nil {
			return nil, err
		} else if ok {
			return v, nil
		} else {
			return nil, ErrReduceNoOutput
		}
	}
}

// MapReduceVoid maps all elements generated from given generate,
// and reduce the output elements with given reducer.
func MapReduceVoid(generate GenerateFunc, mapper MapperFunc, reducer VoidReducerFunc, opts ...Option) error {
	_, err := MapReduce(generate, mapper, func(input <-chan interface{}, writer Writer, cancel func(error)) {
		reducer(input, cancel)
	}, opts...)
	if errors.Is(err, ErrReduceNoOutput) {
		return nil
	}

	return err
}

// WithContext customizes a mapreduce processing accepts a given ctx.
func WithContext(ctx context.Context) Option {
	return func(opts *mapReduceOptions) {
		if ctx != nil {
			opts.ctx = ctx
		}
	}
}

// WithWorkers customizes a mapreduce processing with given workers.
func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func buildSource(generate GenerateFunc, panicChan *onceChan) chan interface{} {
	source := make(chan interface{})
	ants.GoGo(func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			close(source)
		}()

		generate(source)
	})

	return source
}

// drain drains the channel.
func drain(channel <-chan interface{}) {
	// drain the channel
	for range channel {
	}
}

func executeMappers(mCtx mapperContext) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(mCtx.collector)
		drain(mCtx.source)
	}()

	var failed int32
	pool := make(chan mr.PlaceholderType, mCtx.workers) // pool没有关闭？
	defer close(pool)

	writer := newGuardedWriter(mCtx.ctx, mCtx.collector, mCtx.doneChan)
	for atomic.LoadInt32(&failed) == 0 { // 等待failed状态量为true，或者继续循环
		select {
		case <-mCtx.ctx.Done(): // 如果有done，直接退出循环
			return
		case <-mCtx.doneChan:
			return
		case pool <- mr.Placeholder:
			// pool的容量为workers，可以小于并行数，用来控制并发量
			// 如果pool满了，说明到达上限，而且还没执行完毕，等待执行结束的任务
			// 所以workers的数量可以作为参数传入
			item, ok := <-mCtx.source
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			ants.GoGo(func() {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&failed, 1)
						mCtx.panicChan.write(r)
					}
					wg.Done()
					<-pool
				}()

				mCtx.mapper(item, writer)
			})
		}
	}
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		ctx:     context.Background(),
		workers: defaultWorkers,
	}
}

func once(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

type guardedWriter struct {
	ctx     context.Context
	channel chan<- interface{}
	done    <-chan mr.PlaceholderType
}

func newGuardedWriter(ctx context.Context, channel chan<- interface{},
	done <-chan mr.PlaceholderType) guardedWriter {
	return guardedWriter{
		ctx:     ctx,
		channel: channel,
		done:    done,
	}
}

func (gw guardedWriter) Write(v interface{}) {
	select {
	case <-gw.ctx.Done():
		return
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}

type onceChan struct {
	channel chan interface{}
	wrote   int32
}

func (oc *onceChan) write(val interface{}) {
	if atomic.AddInt32(&oc.wrote, 1) > 1 {
		return
	}

	oc.channel <- val
}
