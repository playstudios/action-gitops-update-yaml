"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const core = __importStar(require("@actions/core"));
const yaml_1 = __importDefault(require("yaml"));
const fs_1 = require("fs");
function setValue(str, paths, value) {
    let doc = yaml_1.default.parseDocument(str, {});
    paths.forEach(path => {
        str = doc.setIn(path.split('.'), value);
    });
    return yaml_1.default.stringify(doc);
}
exports.setValue = setValue;
function updateYaml(filePath, keys, value) {
    return __awaiter(this, void 0, void 0, function* () {
        let buffer = yield fs_1.promises.readFile(filePath, {});
        let result = setValue(buffer.toString(), keys, value);
        yield fs_1.promises.writeFile(filePath, result);
    });
}
exports.updateYaml = updateYaml;
function run() {
    return __awaiter(this, void 0, void 0, function* () {
        let workdir = core.getInput("workdir");
        if (workdir.length > 0) {
            process.chdir(workdir);
        }
        let filePath = core.getInput("file_path");
        let keys = core.getInput("key_paths").split(",");
        let value = core.getInput("set_value");
        try {
            yield updateYaml(filePath, keys, value);
        }
        catch (error) {
            core.setFailed(error.message);
        }
    });
}
run().catch((e) => core.setFailed(e.message));
