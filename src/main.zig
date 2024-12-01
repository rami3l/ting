const clap = @import("clap");
const std = @import("std");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const params = comptime clap.parseParamsComptime(
        \\-h, --help             Display this help and exit.
        \\-v, --version          Output version information and exit.
        \\-i, --interval <f32>   Interval between pings, in seconds (default 1)
        \\-n, --count <u16>      Number of tries (default 5)
        \\-p, --port <u16>       Numeric TCP port (default 80)
        \\-w, --timeout <f32>    Maximum time to wait for a response, in seconds (default 5)
        \\<str>...               Hosts to reach
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

    if (res.args.help != 0) {
        const stderr_writer = std.io.getStdErr().writer();
        try stderr_writer.print("ting, yet another TCPing\n\n", .{});
        return clap.help(stderr_writer, clap.Help, &params, .{});
    }
    if (res.args.count) |n| {
        std.debug.print("--count={}\n", .{n});
    }
    for (res.positionals) |pos| {
        std.debug.print("{s}\n", .{pos});
    }
}
