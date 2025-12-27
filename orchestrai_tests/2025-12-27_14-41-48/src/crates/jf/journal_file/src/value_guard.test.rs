#[cfg(test)]
mod value_guard_tests {
    use super::super::*;
    use std::cell::RefCell;
    use std::num::NonZeroU64;
    use std::ops::{Deref, DerefMut};

    // Mock implementation of HashableObject for testing
    #[derive(Debug, Clone)]
    struct MockHashableValue {
        hash_value: u64,
        payload: Vec<u8>,
        next_offset: Option<NonZeroU64>,
    }

    impl MockHashableValue {
        fn new(hash_value: u64, payload: Vec<u8>) -> Self {
            Self {
                hash_value,
                payload,
                next_offset: None,
            }
        }
    }

    impl HashableObject for MockHashableValue {
        fn hash(&self) -> u64 {
            self.hash_value
        }

        fn get_payload(&self) -> &[u8] {
            &self.payload
        }

        fn next_hash_offset(&self) -> Option<NonZeroU64> {
            self.next_offset
        }

        fn object_type() -> ObjectType {
            ObjectType::Data
        }
    }

    impl HashableObjectMut for MockHashableValue {
        fn set_next_hash_offset(&mut self, offset: NonZeroU64) {
            self.next_offset = Some(offset);
        }

        fn set_payload(&mut self, data: &[u8]) {
            self.payload = data.to_vec();
        }
    }

    #[test]
    fn test_value_guard_new_creates_instance() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![1, 2, 3]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);

        assert_eq!(guard.offset.get(), 100);
    }

    #[test]
    fn test_value_guard_offset_returns_correct_value() {
        let offset = NonZeroU64::new(12345).unwrap();
        let value = MockHashableValue::new(99, vec![]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        assert_eq!(guard.offset(), offset);
    }

    #[test]
    fn test_value_guard_offset_with_zero_max_values() {
        let offset = NonZeroU64::new(u64::MAX).unwrap();
        let value = MockHashableValue::new(0, vec![]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        assert_eq!(guard.offset().get(), u64::MAX);
    }

    #[test]
    fn test_value_guard_deref_returns_reference() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![1, 2, 3]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let deref_value: &MockHashableValue = guard.deref();

        assert_eq!(deref_value.hash(), 42);
        assert_eq!(deref_value.get_payload(), &[1, 2, 3]);
    }

    #[test]
    fn test_value_guard_deref_multiple_calls() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(55, vec![10, 20]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let ref1 = guard.deref();
        let ref2 = guard.deref();

        assert_eq!(ref1.hash(), ref2.hash());
    }

    #[test]
    fn test_value_guard_deref_mut_returns_mutable_reference() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![1, 2, 3]);
        let in_use = RefCell::new(true);

        let mut guard = ValueGuard::new(offset, value, &in_use);
        let deref_mut_value: &mut MockHashableValue = guard.deref_mut();

        deref_mut_value.hash_value = 99;
        assert_eq!(guard.hash_value, 99);
    }

    #[test]
    fn test_value_guard_deref_mut_allows_payload_modification() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![1, 2, 3]);
        let in_use = RefCell::new(true);

        let mut guard = ValueGuard::new(offset, value, &in_use);
        {
            let deref_mut_value = guard.deref_mut();
            deref_mut_value.payload = vec![4, 5, 6];
        }

        assert_eq!(guard.payload, vec![4, 5, 6]);
    }

    #[test]
    fn test_value_guard_drop_sets_in_use_flag_to_false() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        {
            let _guard = ValueGuard::new(offset, value, &in_use);
            assert_eq!(*in_use.borrow(), true);
        }

        assert_eq!(*in_use.borrow(), false);
    }

    #[test]
    fn test_value_guard_drop_flag_already_false() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(false);

        {
            let _guard = ValueGuard::new(offset, value, &in_use);
            assert_eq!(*in_use.borrow(), false);
        }

        assert_eq!(*in_use.borrow(), false);
    }

    #[test]
    fn test_value_guard_hashable_object_hash_delegation() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(12345, vec![]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let result = HashableObject::hash(&guard);

        assert_eq!(result, 12345);
    }

    #[test]
    fn test_value_guard_hashable_object_get_payload_delegation() {
        let offset = NonZeroU64::new(100).unwrap();
        let payload = vec![7, 8, 9];
        let value = MockHashableValue::new(42, payload.clone());
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let result = guard.get_payload();

        assert_eq!(result, &payload[..]);
    }

    #[test]
    fn test_value_guard_hashable_object_next_hash_offset_none() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let result = guard.next_hash_offset();

        assert_eq!(result, None);
    }

    #[test]
    fn test_value_guard_hashable_object_next_hash_offset_some() {
        let offset = NonZeroU64::new(100).unwrap();
        let mut value = MockHashableValue::new(42, vec![]);
        value.next_offset = NonZeroU64::new(200);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let result = guard.next_hash_offset();

        assert_eq!(result, NonZeroU64::new(200));
    }

    #[test]
    fn test_value_guard_hashable_object_type() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        let _guard = ValueGuard::new(offset, value, &in_use);
        let object_type = MockHashableValue::object_type();

        assert_eq!(object_type, ObjectType::Data);
    }

    #[test]
    fn test_value_guard_hashable_object_mut_set_next_hash_offset() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        let mut guard = ValueGuard::new(offset, value, &in_use);
        let new_offset = NonZeroU64::new(500).unwrap();

        HashableObjectMut::set_next_hash_offset(&mut guard, new_offset);

        assert_eq!(guard.next_hash_offset(), Some(new_offset));
    }

    #[test]
    fn test_value_guard_hashable_object_mut_set_payload() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![1, 2, 3]);
        let in_use = RefCell::new(true);

        let mut guard = ValueGuard::new(offset, value, &in_use);
        let new_payload = [4, 5, 6];

        HashableObjectMut::set_payload(&mut guard, &new_payload);

        assert_eq!(guard.get_payload(), &new_payload);
    }

    #[test]
    fn test_value_guard_debug_impl() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![1, 2, 3]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let debug_string = format!("{:?}", guard);

        assert!(!debug_string.is_empty());
        assert!(debug_string.contains("ValueGuard"));
    }

    #[test]
    fn test_value_guard_with_empty_payload() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);

        assert_eq!(guard.get_payload().len(), 0);
    }

    #[test]
    fn test_value_guard_with_large_payload() {
        let offset = NonZeroU64::new(100).unwrap();
        let large_payload = vec![42; 10000];
        let value = MockHashableValue::new(99, large_payload.clone());
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);

        assert_eq!(guard.get_payload().len(), 10000);
        assert_eq!(guard.get_payload(), &large_payload[..]);
    }

    #[test]
    fn test_value_guard_drop_multiple_times_in_sequence() {
        let in_use1 = RefCell::new(true);
        let in_use2 = RefCell::new(true);

        {
            let offset = NonZeroU64::new(100).unwrap();
            let value = MockHashableValue::new(42, vec![]);
            let _guard1 = ValueGuard::new(offset, &in_use1);
        }

        {
            let offset = NonZeroU64::new(200).unwrap();
            let value = MockHashableValue::new(99, vec![]);
            let _guard2 = ValueGuard::new(offset, &in_use2);
        }

        assert_eq!(*in_use1.borrow(), false);
        assert_eq!(*in_use2.borrow(), false);
    }

    #[test]
    fn test_value_guard_generic_with_different_types() {
        #[derive(Debug)]
        struct SimpleValue(u32);

        impl HashableObject for SimpleValue {
            fn hash(&self) -> u64 {
                self.0 as u64
            }
            fn get_payload(&self) -> &[u8] {
                &[]
            }
            fn next_hash_offset(&self) -> Option<NonZeroU64> {
                None
            }
            fn object_type() -> ObjectType {
                ObjectType::Entry
            }
        }

        impl HashableObjectMut for SimpleValue {
            fn set_next_hash_offset(&mut self, _offset: NonZeroU64) {}
            fn set_payload(&mut self, _data: &[u8]) {}
        }

        let offset = NonZeroU64::new(100).unwrap();
        let value = SimpleValue(123);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);

        assert_eq!(guard.hash(), 123);
        assert_eq!(*in_use.borrow(), true);
    }

    #[test]
    fn test_value_guard_offset_getter_is_nonzero() {
        let offset = NonZeroU64::new(1).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        let guard = ValueGuard::new(offset, value, &in_use);
        let retrieved_offset = guard.offset();

        // NonZeroU64 is guaranteed to be non-zero
        assert_eq!(retrieved_offset.get(), 1);
        assert_ne!(retrieved_offset.get(), 0);
    }

    #[test]
    fn test_value_guard_preserves_offset_immutably() {
        let offset = NonZeroU64::new(54321).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        let mut guard = ValueGuard::new(offset, value, &in_use);
        let offset1 = guard.offset();

        // Even after mutation
        let _ = guard.deref_mut();
        let offset2 = guard.offset();

        assert_eq!(offset1, offset2);
    }

    #[test]
    fn test_value_guard_set_payload_empty_to_nonempty() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![]);
        let in_use = RefCell::new(true);

        let mut guard = ValueGuard::new(offset, value, &in_use);
        assert_eq!(guard.get_payload().len(), 0);

        HashableObjectMut::set_payload(&mut guard, &[1, 2, 3, 4, 5]);
        assert_eq!(guard.get_payload(), &[1, 2, 3, 4, 5]);
    }

    #[test]
    fn test_value_guard_set_payload_nonempty_to_empty() {
        let offset = NonZeroU64::new(100).unwrap();
        let value = MockHashableValue::new(42, vec![1, 2, 3]);
        let in_use = RefCell::new(true);

        let mut guard = ValueGuard::new(offset, value, &in_use);
        assert_eq!(guard.get_payload().len(), 3);

        HashableObjectMut::set_payload(&mut guard, &[]);
        assert_eq!(guard.get_payload().len(), 0);
    }
}