const fs = require('fs');

async function parse() {
  return [0.0, 0.0, 0];
}

function compFloats(f1, f2) {
  const precision = 10;
  const int1 = Math.floor(f1 * precision + 0.5);
  const int2 = Math.floor(f2 * precision + 0.5);
  return int1 === int2;
}

async function main() {
  const tp1 = performance.now();
  const res = await parse();
  const tp2 = performance.now();
  const execTime = tp2 - tp1;

  const data = fs.readFileSync('points-verify.txt', 'utf8');
  const parts = data.slice(0, -1).split(',');
  const p1 = parseFloat(parts[0]);
  const p2 = parseFloat(parts[1]);
  const lines = parseInt(parts[2], 10);

  if (res[2] !== lines) {
    throw new Error(`Excpected number of lines to be ${lines}, got: ${res[2]}`);
  }

  if (!compFloats(res[0], p1)) {
    throw new Error(`Excpected first number to be ${p1}, got: ${res[0]}`);
  }

  if (!compFloats(res[1], p2)) {
    throw new Error(`Excpected first number to be ${p2}, got: ${res[1]}`);
  }

  return execTime;
}

(async () => {
  let bestTime = 100000.0
  while (true) {
    const execTime = (await main()) / 1000.0;
    if (execTime < bestTime) {
      bestTime = execTime;
      console.log(`Execution time: ${bestTime}`);
    }
  }
})();
