package context

import (
	"context"
	"testing"
	"time"
)

type mykey struct {
}

func TestContext_WithValue(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, mykey{}, "haha")
	val := ctx.Value(mykey{}).(string)
	t.Log(val)
	newVal := ctx.Value("不存在的key")
	val, ok := newVal.(string)
	if !ok {
		t.Log("类型不对")
		return
	}
	t.Log(val)

}

func TestContext_WithCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	<-ctx.Done()
	t.Log("hellw cancel", ctx.Err())
}

// 超时控制
func TestContext_WithDeadline(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	deadline, _ := ctx.Deadline()
	t.Log("deadline:", deadline)
	defer cancel()

	<-ctx.Done()
	t.Log("deadline", ctx.Err())
}

// 超时控制
func TestContext_WithTimeout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	deadline, _ := ctx.Deadline()
	t.Log("deadline:", deadline)
	defer cancel()

	<-ctx.Done()
	t.Log("deadline", ctx.Err())
}

func TestContext_Parent(t *testing.T) {
	ctx := context.Background()
	parent := context.WithValue(ctx, "my-key", "my_value")
	child := context.WithValue(ctx, "my-key", "my new value")

	t.Log("parent my-key", parent.Value("my-key"))
	t.Log("child my-key", child.Value("my-key"))

	child2, cancel := context.WithTimeout(parent, time.Second)
	defer cancel()
	t.Log("child2 my-key", child2.Value("my-key"))

	child3 := context.WithValue(parent, "new-key", "child3 value")

	t.Log("parent new-key", parent.Value("new-key"))
	t.Log("child3 new-key", child3.Value("new-key"))

	// 父节点可拿到子节点的方案，避免用，
	parent1 := context.WithValue(parent, "map", map[string]string{})
	child4, cancel := context.WithTimeout(parent1, time.Second)
	defer cancel()
	m := child4.Value("map").(map[string]string)
	m["key1"] = "value1"
	nm := parent1.Value("map").(map[string]string)
	t.Log("parent key1", nm["key1"])
}
