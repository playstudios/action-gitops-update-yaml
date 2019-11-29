// import * as process from 'process';
// import * as cp from 'child_process';
// import * as path from 'path';

import {setValue} from '../src/main';

test('setValue', () => {
    const tests = [
        {
            input: 'a:\n  b: hello',
            paths: ['a.b'],
            value: 'world',
            output: 'a:\n  b: world'
        },
        {
            input: 'a:\n  b: hello\nc:\n  d: abc',
            paths: ['a.b', 'c.d'],
            value: 'world',
            output: 'a:\n  b: world\nc:\n  d: world'
        },

        // Currently unsupported. https://github.com/eemeli/yaml/issues/131
        // {
        //     input: 'a:\n  b: hello',
        //     paths: ['a.b.c'],
        //     value: 'world',
        //     output: 'a:\n  b: world'
        // }
    ];

    tests.forEach(test => {
        expect(setValue(test.input, test.paths, test.value)).toMatch(test.output);
    });
});

// test('wait 500 ms', async() => {
//     const start = new Date();
//     await wait(500);
//     const end = new Date();
//     var delta = Math.abs(end.getTime() - start.getTime());
//     expect(delta).toBeGreaterThan(450);
// });
//
// // shows how the runner will run a javascript action with env / stdout protocol
// test('test runs', () => {
//     process.env['INPUT_MILLISECONDS'] = '500';
//     const ip = path.join(__dirname, '..', 'lib', 'main.js');
//     const options: cp.ExecSyncOptions = {
//         env: process.env
//     };
//     console.log(cp.execSync(`node ${ip}`, options).toString());
// });
