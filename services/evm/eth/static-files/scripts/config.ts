import fs from "fs";
const { PWD } = process.env;
const DEPLOYMENTS_PATH = `${PWD}/deployments.json`;

export class Deployments {
  map: Map<string, any>;
  private static instance: Deployments;

  constructor(map: any) {
    this.map = map;
  }

  public get(target: string) {
    return this.map.get(target);
  }

  public set(target: string, data: any) {
    this.map.set(target, data);
  }

  public static getDefault() {
    if (!this.instance) {
      const data = fs.readFileSync(DEPLOYMENTS_PATH);
      const json = JSON.parse(data.toString());
      const map = new Map(Object.entries(json));
      this.instance = new this(map);
    }
    return this.instance;
  }

  public save() {
    fs.writeFileSync(
      DEPLOYMENTS_PATH,
      JSON.stringify(Object.fromEntries(this.map)),
      "utf-8"
    );
  }
}
