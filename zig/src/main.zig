const std = @import("std");
const dbg = std.debug;

const Result = struct { min: f64, max: f64, lines: u64 };

fn parse() !Result {
    return Result{ .min = 0.0, .max = 0.0, .lines = 0 };
}

fn compFloats(f1: f64, f2: f64) bool {
    const precision = 10;
    const int1: i64 = @intFromFloat(f1 * precision + 0.5);
    const int2: i64 = @intFromFloat(f2 * precision + 0.5);
    return int1 == int2;
}

pub fn main() !void {
    var bestTime: u64 = 999_999_999_999;

    var file = try std.fs.cwd().openFile("points-verify.txt", .{});
    var buf: [1024]u8 = undefined;
    const bytesRead = try file.readAll(&buf);
    const line = buf[0 .. bytesRead - 1];
    var parts = std.mem.splitScalar(u8, line, ',');

    const f1 = try std.fmt.parseFloat(f64, parts.first());
    const f2 = try std.fmt.parseFloat(f64, parts.next().?);
    const lines = try std.fmt.parseInt(u64, parts.rest(), 10);

    while (true) {
        const tp1 = try std.time.Instant.now();
        const res = try parse();
        const tp2 = try std.time.Instant.now();
        const elapsed = tp2.since(tp1);

        if (lines != res.lines) {
            dbg.panic("Expected number of lines to be {d} got {d}\n", .{ lines, res.lines });
        }

        if (!compFloats(f1, res.min)) {
            dbg.panic("Expected first number to be {d} got {d}\n", .{ f1, res.min });
        }

        if (!compFloats(f2, res.max)) {
            // if (!std.mem.eql(u8, rmax, max)) {
            dbg.panic("Expected second number to be {d} got {d}\n", .{ f2, res.max });
        }
        if (elapsed < bestTime) {
            dbg.print("Execution time: {s}\n", .{std.fmt.fmtDuration(elapsed)});
            bestTime = elapsed;
        }
    }
}
