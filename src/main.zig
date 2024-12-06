const std = @import("std");
const fmt = std.fmt;

const clap = @import("clap");
const ting = @import("ting");

const meta = @import("./meta.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const params = comptime clap.parseParamsComptime(blk: {
        const t = ting.Tcping{ .host = "" };
        break :blk fmt.comptimePrint(
            \\-h, --help             Display this help and exit.
            \\-v, --version          Output version information and exit.
            \\-i, --interval <f32>   Interval between pings, in seconds (default: {d:.1})
            \\-c, --count <u16>      Number of tries (default: {s})
            \\-p, --port <u16>       Numeric TCP port (default: {d})
            \\-w, --timeout <f32>    Maximum time to wait for a response, in seconds (default: {d:.1})
            \\<str>                  Host to reach
        , .{
            t.interval_s,
            if (t.count) |c| fmt.comptimePrint("{d}", c) else "unlimited",
            t.port,
            t.timeout_s,
        });
    });

    var diag = clap.Diagnostic{};
    const parsed = clap.parse(clap.Help, &params, clap.parsers.default, .{
        .diagnostic = &diag,
        .allocator = gpa.allocator(),
    }) catch |e| {
        // Report useful error and exit.
        diag.report(std.io.getStdErr().writer(), e) catch {};
        return e;
    };
    defer parsed.deinit();

    const stderr_writer = std.io.getStdErr().writer();

    if (parsed.args.help != 0) {
        try stderr_writer.print("ting, yet another TCPing\n", .{});
        return clap.help(stderr_writer, clap.Help, &params, .{
            .spacing_between_parameters = 0,
        });
    }

    if (parsed.args.version != 0) {
        try stderr_writer.print("ting v{} [zig v{}]", .{
            meta.version,
            @import("builtin").zig_version,
        });
        return;
    }

    const count_positionals = parsed.positionals.len;
    if (count_positionals != 1) {
        try stderr_writer.print("error: expected a single host (found {})\n", .{count_positionals});
        std.process.exit(1);
    }

    var t = ting.Tcping{
        .host = parsed.positionals[0],
        .count = parsed.args.count,
    };
    if (parsed.args.interval) |it| t.interval_s = it;
    if (parsed.args.count) |it| t.count = it;
    if (parsed.args.port) |it| t.port = it;
    if (parsed.args.timeout) |it| t.timeout_s = it;

    const stdout_writer = std.io.getStdOut().writer();

    const durations = try t.ping(gpa.allocator(), stdout_writer);
    defer durations.deinit();
    try t.report(durations.items, stdout_writer);
}
