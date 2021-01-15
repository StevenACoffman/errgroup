package errgroup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type errgroupSuite struct{ suite.Suite }

func (suite *errgroupSuite) TestPanicWithString() {
	g, ctx := WithContext(context.Background())
	g.Go(func() error { panic("oh noes") })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "oh noes")

	// ctx should now be canceled.
	suite.Require().Error(ctx.Err())
}

func (suite *errgroupSuite) TestPanicWithError() {
	g, ctx := WithContext(context.Background())

	panicErr := errors.New("oh noes")
	g.Go(func() error { panic(panicErr) })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "oh noes")
	suite.Require().True(errors.Is(err, panicErr))

	// ctx should now be canceled.
	suite.Require().Error(ctx.Err())
}

func (suite *errgroupSuite) TestPanicWithOtherValue() {
	g, ctx := WithContext(context.Background())

	panicVal := struct {
		int
		string
	}{1234567890, "oh noes"}
	g.Go(func() error { panic(panicVal) })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "oh noes")
	suite.Require().Contains(err.Error(), "1234567890")

	// ctx should now be canceled.
	suite.Require().Error(ctx.Err())
}

func (suite *errgroupSuite) TestError() {
	g, ctx := WithContext(context.Background())

	goroutineErr := errors.New("oh noes")
	g.Go(func() error { return goroutineErr })
	// this function ensures that the panic in fact cancels the context, by not
	// returning until it's been cancelled; it should return context.Canceled
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	// Wait() will finish only once all goroutines do, but returns the first
	// error
	err := g.Wait()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "oh noes")
	suite.Require().True(errors.Is(err, goroutineErr))

	// ctx should now be canceled.
	suite.Require().Error(ctx.Err())
}

func (suite *errgroupSuite) TestSuccess() {
	g, ctx := WithContext(context.Background())

	g.Go(func() error { return nil })
	// since no goroutine errored, ctx.Err() should be nil
	// (until all goroutines are done)
	g.Go(ctx.Err)

	err := g.Wait()
	suite.Require().NoError(err)

	// ctx should now still be canceled.
	suite.Require().Error(ctx.Err())
}

func (suite *errgroupSuite) TestManyGoroutines() {
	n := 100
	g, ctx := WithContext(context.Background())

	for i := 0; i < n; i++ {
		// put in a bunch of goroutines that just return right away
		g.Go(func() error { return nil })
		// and also a bunch that wait for the error
		g.Go(func() error {
			<-ctx.Done()
			return ctx.Err()
		})
	}

	// finally, put in a panic
	g.Go(func() error { panic("oh noes") })

	// as before, Wait() will finish only once all goroutines do, but returns
	// the first error (namely the panic)
	err := g.Wait()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "oh noes")

	// ctx should now be canceled.
	suite.Require().Error(ctx.Err())
}

func (suite *errgroupSuite) TestZeroGroupPanic() {
	var g Group

	// either of these could happen first, since a zero group does not cancel
	g.Go(func() error { panic("oh noes") })
	g.Go(func() error { return nil })

	// Wait() still returns the error.
	err := g.Wait()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "oh noes")
}

func (suite *errgroupSuite) TestZeroGroupSuccess() {
	var g Group

	g.Go(func() error { return nil })
	g.Go(func() error { return nil })

	err := g.Wait()
	suite.Require().NoError(err)
}

func TestErrgroup(t *testing.T) {
	suite.Run(t, new(errgroupSuite))
}