import fs from 'fs';
const {PWD} = process.env
const DEPLOYMENTS_PATH = `${PWD}/deployments.json`
const CHAIN_CONFIG_PATH = `${PWD}/chain_config.json`

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

  public getSrc() {
    return this.map.get('link').src;
  }

  public getDst() {
    return this.map.get('link').dst;
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
    fs.writeFileSync(DEPLOYMENTS_PATH, JSON.stringify(Object.fromEntries(this.map)), 'utf-8')
  }
}

export class ChainConfig {
  private static map: Map<String, any>;

  public static getProp(key: string) {
    if (!this.map) {
      const data = fs.readFileSync(CHAIN_CONFIG_PATH);
      const json = JSON.parse(data.toString());
      this.map = new Map(Object.entries(json));
    }
    return this.map.get(key);
  }

  public static getChain(target: string) {
    const chains = this.getProp('chains');
    const config = new Map(Object.entries(chains));
    return config.get(target);
  }

  public static getLink() {
    return this.getProp('link');
  }
}

export function chainType(chain: any) {
  return chain.network.split(".")[1];
}
