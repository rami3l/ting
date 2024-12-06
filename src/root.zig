const std = @import("std");
const mem = std.mem;
const meta = std.meta;
const net = std.net;
const posix = std.posix;
const sc = std.c;
const time = std.time;

const ic = @cImport({
    @cInclude("fcntl.h");
    @cInclude("sys/select.h");
});

const ArrayList = std.ArrayList;

pub const Tcping = struct {
    count: ?u16 = null,
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
            const sock = try posix.socket(sc.AF.INET, sc.SOCK.STREAM, 0);
            defer posix.close(sock);

            const flags = try posix.fcntl(sock, sc.F.GETFL, 0);
            _ = try posix.fcntl(sock, sc.F.SETFL, flags | ic.O_NONBLOCK);

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
                else => break error.CouldNotConnect,
            }

            const end = try time.Instant.now();
            break .{ end.since(start), addr };
        } else error.CouldNotConnect;
    }

    pub fn ping(self: *const Self, alloc: mem.Allocator, log_writer: anytype) !ArrayList(?u64) {
        var durations = ArrayList(?u64).init(alloc);
        try log_writer.print("TCPING {s}:{d}\n", .{ self.host, self.port });
        var i: meta.Child(@TypeOf(self.count)) = 0;
        return loop: while (if (self.count) |c| i < c else true) : (i += 1) {
            const duration, const addr = self.probe(alloc, log_writer) catch |e| {
                try log_writer.print("error: ", .{});
                switch (e) {
                    error.UnknownHostName => {
                        try log_writer.print("cannot resolve host {s}\n", .{self.host});
                        break :loop durations;
                    },
                    error.CouldNotConnect => {
                        try log_writer.print("connection timed out after {d:.1}s for seq={d}\n", .{
                            self.timeout_s,
                            i,
                        });
                        try durations.append(null);
                        continue :loop;
                    },
                    else => break :loop e,
                }
            };
            try log_writer.print("addr={} seq={d} time={d:.3}ms\n", .{ addr, i, ms_from_ns(duration) });
            try durations.append(duration);
            time.sleep(@intFromFloat(self.interval_s * @as(f32, @floatFromInt(time.ns_per_s))));
        } else durations;
    }

    pub fn report(self: *const Self, durations: []?u64, log_writer: anytype) !void {
        const count = durations.len;
        if (count == 0) return;

        var oks: u64 = 0;
        var min: u64, var max: u64 = .{ std.math.maxInt(u64), 0 };
        var total: u128 = 0;
        for (durations) |duration| if (duration) |d| {
            oks += 1;
            min = @min(min, d);
            max = @max(max, d);
            total += d;
        };

        try log_writer.print("--- {s} tcping statistics ---\n", .{self.host});
        try log_writer.print("{d} connections, {d} succeeded, {d} failed, {d:.1}% success rate\n", .{
            count,
            oks,
            count - oks,
            @as(f32, @floatFromInt(oks)) / @as(f32, @floatFromInt(count)) * 100.0,
        });
        if (oks == 0) return;
        try log_writer.print("minimum = {d:.3}ms, maximum = {d:.3}ms, average = {d:.3}ms\n", .{
            ms_from_ns(min),
            ms_from_ns(max),
            ms_from_ns(total) / @as(f32, @floatFromInt(oks)),
        });
    }

    pub const Error = error{
        CouldNotConnect,
        UnknownHostName,
        Interrupted,
    };
};

fn ms_from_ns(ns: anytype) f32 {
    const ms_per_ns = 1 / @as(f32, @floatFromInt(time.ns_per_ms));
    return @as(f32, @floatFromInt(ns)) * ms_per_ns;
}
