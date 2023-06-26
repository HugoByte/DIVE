import fs from 'fs';
import path from 'path';

export class Jar {
  public static readFromFile(base: string | undefined, project: string, version: string = '0.1.0') {
    if (!base) {
      base = "../javascore";
    }
    const build = "build/libs";
    const name = project.replace("/", "-");
    const optJar = `${name}-${version}-optimized.jar`;
    const fullPath = path.join(base, project, build, optJar);
    return fs.readFileSync(fullPath).toString('hex')
  }
}
