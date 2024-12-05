const clap = @import("clap");
const std = @import("std");

const ting = @import("ting");

const meta = @import("./meta.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const params = comptime clap.parseParamsComptime(
        \\-h, --help             Display this help and exit.
        \\-v, --version          Output version information and exit.
        \\-i, --interval <f32>   Interval between pings, in seconds (default 1)
        \\-c, --count <u16>      Number of tries (default 5)
        \\-p, --port <u16>       Numeric TCP port (default 80)
        \\-w, --timeout <f32>    Maximum time to wait for a response, in seconds (default 5)
        \\<str>                  Host to reach
    );

    var diag = clap.Diagnostic{};
    var res = clap.parse(clap.Help, &params, clap.parsers.default, .{
        .diagnostic = &diag,
        .allocator = gpa.allocator(),
    }) catch |err| {
        // Report useful error and exit.
        diag.report(std.io.getStdErr().writer(), err) catch {};
        return err;
    };
    defer res.deinit();

    const stderr_writer = std.io.getStdErr().writer();

    if (res.args.help != 0) {
        try stderr_writer.print("ting, yet another TCPing\n\n", .{});
        return clap.help(stderr_writer, clap.Help, &params, .{});
    }

    if (res.args.version != 0) {
        try stderr_writer.print("ting v{} [zig v{}]", .{
            meta.version,
            @import("builtin").zig_version,
        });
        return;
    }

    const count_positionals = res.positionals.len;
    if (count_positionals != 1) {
        try stderr_writer.print("error: expected a single host (found {})\n", .{count_positionals});
        std.process.exit(1);
    }
    const host = res.positionals[0];

    var cfg = ting.Tcping{
        .host = host,
    };
    if (res.args.interval) |it| cfg.interval_s = it;
    if (res.args.count) |it| cfg.count = it;
    if (res.args.port) |it| cfg.port = it;
    if (res.args.timeout) |it| cfg.timeout_s = it;

    try cfg.ping(gpa.allocator(), stderr_writer);
}
