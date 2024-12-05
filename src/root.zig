const std = @import("std");
const c = std.c;
const mem = std.mem;
const net = std.net;
const posix = std.posix;
const time = std.time;

const ic = @cImport({
    @cInclude("fcntl.h");
    @cInclude("sys/select.h");
});

pub const Tcping = struct {
    count: u16 = 5,
    port: u16 = 80,
    interval_s: f32 = 1.0,
    timeout_s: f32 = 5.0,
    host: []const u8,

    const Self = @This();

    /// Returns the time required to connect to the host via TCP (in nanoseconds).
    pub fn probe(self: *const Self, alloc: mem.Allocator, log_writer: anytype) !struct { u64, net.Address } {
        try log_writer.print("connecting... ", .{});
        const timeout_ns: u64 = @intFromFloat(self.timeout_s * @as(f32, @floatFromInt(time.ns_per_s)));

        const addr_list = try net.getAddressList(alloc, self.host, self.port);
        defer addr_list.deinit();

        return for (addr_list.addrs) |addr| {
            // https://stackoverflow.com/a/2597774
            const sock = try posix.socket(c.AF.INET, c.SOCK.STREAM, 0);
            defer posix.close(sock);

            const flags = try posix.fcntl(sock, c.F.GETFL, 0);
            _ = try posix.fcntl(sock, c.F.SETFL, flags | ic.O_NONBLOCK);

            const start = try time.Instant.now();
            while (true) {
                posix.connect(sock, &addr.any, @sizeOf(@TypeOf(addr.any))) catch |e| switch (e) {
                    error.WouldBlock => continue,
                    error.ConnectionPending => break,
                    else => return e,
                };
                break;
            }

            var fd_set = ic.fd_set{};
            ic.FD_SET(sock, &fd_set);

            var timeout = posix.timeval{
                .tv_sec = @intCast(timeout_ns / time.ns_per_s),
                .tv_usec = @intCast(timeout_ns % time.ns_per_s),
            };
            switch (ic.select(sock + 1, null, &fd_set, null, @ptrCast(&timeout))) {
                1 => try posix.getsockoptError(sock),
                else => return error.CouldNotConnect,
            }

            const end = try time.Instant.now();
            try log_writer.print("addr={} ", .{addr});
            break .{ end.since(start), addr };
        } else error.CouldNotConnect;
    }

    pub fn ping(self: *const Self, alloc: mem.Allocator, log_writer: anytype) !void {
        try log_writer.print("TCPING {s}:{d} with a timeout of {d:.1}s:\n", .{ self.host, self.port, self.timeout_s });
        return loop: for (0..self.count) |i| {
            const res = self.probe(alloc, log_writer) catch |e| {
                try log_writer.print("error: ", .{});
                switch (e) {
                    error.UnknownHostName => {
                        try log_writer.print("cannot resolve host {s}\n", .{self.host});
                        break :loop;
                    },
                    error.CouldNotConnect => {
                        try log_writer.print("cannot connect to {s}:{d} for seq={d}\n", .{ self.host, self.port, i });
                        continue :loop;
                    },
                    else => break :loop e,
                }
            };
            try log_writer.print("seq={d} time={}ns\n", .{ i, res[0] });
            time.sleep(@intFromFloat(self.interval_s * @as(f32, @floatFromInt(time.ns_per_s))));
        };
    }

    pub const Error = error{
        CouldNotConnect,
        UnknownHostName,
    };
};

// const testing = std.testing;
//
// export fn add(a: i32, b: i32) i32 {
//     return a + b;
// }
//
// test "basic add functionality" {
//     try testing.expect(add(3, 7) == 10);
// }
