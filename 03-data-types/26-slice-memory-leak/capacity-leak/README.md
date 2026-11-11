# Slice Memory Leak (Capacity Leak)

This example demonstrates a classic Go memory leak caused by sharing a backing array between slices.

## What's happening

1. `receiveMessage()` allocates a large buffer — a 1,000,000-byte slice.
2. `getMessageType()` returns `msg[:5]`, a reslice of just the first 5 bytes.
3. Slicing **does not copy data**. The returned 5-byte slice shares the same underlying array as the original 1 MB message.
4. `storeMessageType()` keeps that 5-byte slice alive. Because the slice header still points into the 1 MB backing array, the garbage collector **cannot free the entire 1 MB** — it is pinned by the tiny slice.

So a 5-byte value ends up retaining a full megabyte of memory. In a long-running loop like `consumeMessages()`, this causes memory to grow without bound even though you're only "using" a few bytes per message.

## Key concept

A Go slice header holds a pointer, length, and capacity — all pointing to a shared backing array. Any slice derived from that array keeps the **whole array** alive, not just the visible portion. When you hold onto a small reslice of a large buffer, you leak the entire buffer's capacity.

## The fix

`getMessageTypeWithCopy()` shows the correct pattern:

```go
func getMessageTypeWithCopy(msg []byte) []byte {
	msgType := make([]byte, 5)
	copy(msgType, msg)
	return msgType
}
```

By allocating a new slice and copying only the 5 bytes you need, you break the shared backing array. The original 1 MB message buffer can then be garbage collected normally.

## Rule of thumb

If you extract a small slice from a large buffer and intend to keep it around beyond the lifetime of the original buffer, copy the data you need instead of reslicing.

## Running the example

`main()` drives `consumeMessages(n)` with both extraction functions, running the **fixed** path first, then the **leaky** path. `storeMessageType()` retains each result so the leak is real:

```
--- Leaky: getMessageType (msg[:5]) ---
Alloc before: 350 KB
Alloc after 100 messages: 98751 KB
Alloc after GC (still leaked): 98754 KB

--- Fixed: getMessageTypeWithCopy ---
Alloc before: 297 KB
Alloc after 100 messages: 1331 KB
Alloc after GC: 350 KB
```

With the copy, only the small 5-byte slices are retained and the 1 MB buffers are freed by GC. With the reslice (`msg[:5]`), each retained 5-byte slice pins its entire 1 MB backing array, so GC cannot reclaim ~96 MB even after a full collection.
