// Copyright 2021-2025, Offchain Labs, Inc.
// For license information, see https://github.com/OffchainLabs/nitro/blob/master/LICENSE

use std::io::Write;

pub fn write_wat_ops(wat: &mut Vec<u8>, number_of_ops_per_loop_iteration: usize) {
    for _ in 0..number_of_ops_per_loop_iteration {
        wat.write_all(b"            (block\n").unwrap();
        wat.write_all(b"                (block \n").unwrap();
        wat.write_all(b"                    (block \n").unwrap();
        wat.write_all(b"                        br 2\n").unwrap(); // it will jump to the end of the last block
        wat.write_all(b"                    )\n").unwrap();
        wat.write_all(b"                )\n").unwrap();
        wat.write_all(b"            )\n").unwrap();
    }
}
