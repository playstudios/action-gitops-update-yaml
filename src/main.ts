import * as core from '@actions/core';
import yaml from 'yaml';
import {promises as fs} from 'fs';

export function setValue(str: string, paths: string[], value: string): string {
    let doc = yaml.parseDocument(str, {});
    paths.forEach(path => {
        str = doc.setIn(path.split('.'), value);
    });
    return yaml.stringify(doc);
}

export async function updateYaml(filePath: string, keys: string[], value: string): Promise<void> {
    let buffer = await fs.readFile(filePath, {});
    let result = setValue(buffer.toString("UTF-8"), keys, value);
    await fs.writeFile(filePath, result, {encoding: "UTF-8"})
}

async function run() {
    let filePath = core.getInput("file_path", {required: true});
    let keys = core.getInput("keys", {required: true}).split(",");
    let value = core.getInput("set_value", {required: true});
    try {
        await updateYaml(filePath, keys, value);
    } catch (error) {
        core.setFailed(error.message);
    }
}

run().catch((e) => core.setFailed(e.message));
