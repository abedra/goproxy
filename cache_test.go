package main

import "testing"

func TestLifecycle(t *testing.T) {
	lru := NewCache(0)
	lru.Add("1.1.1.1", "value")
	lru.Add("2.2.2.2", "oldest")

	val, ok := lru.Get("1.1.1.1")
	if !ok {
		t.Fatalf("cache hit = %v; want %v", ok, !ok)
	} else if ok && val != "value" {
		t.Fatalf("Expected Get to return 'value' but go %v", val)
	}

	lru.Remove("1.1.1.1")
	if _, ok := lru.Get("1.1.1.1"); ok {
		t.Fatal("Remove failed")
	}

	lru.RemoveOldest()
	if lru.Len() != 0 {
		t.Fatalf("Expected size to be 0, got %v", lru.Len())
	}
}

func TestEviction(t *testing.T) {
	lru := NewCache(2)
	lru.Add("1.1.1.1", "first")
	lru.Add("2.2.2.2", "second")
	lru.Add("3.3.3.3", "third")

	if lru.Len() > 2 {
		t.Fatalf("Expected size to be 2, but got %v", lru.Len())
	}

	if _, ok := lru.Get("1.1.1.1"); ok {
		t.Fatal("Expected 1.1.1.1 to be removed by LRU eviction")
	}

	if _, ok := lru.Get("2.2.2.2"); !ok {
		t.Fatal("Expected 2.2.2.2 to be present")
	}

	if _, ok := lru.Get("3.3.3.3"); !ok {
		t.Fatal("Expected 3.3.3.3 to be present")
	}

	lru.Get("2.2.2.2")
	lru.Add("4.4.4.4", "fourth")

	if _, ok := lru.Get("3.3.3.3"); ok {
		t.Fatal("Expected 3.3.3.3 to be evicted")
	}

	if _, ok := lru.Get("2.2.2.2"); !ok {
		t.Fatal("Expected 2.2.2.2 to be present")
	}
}
