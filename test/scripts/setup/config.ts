import fs from 'fs';
const {PWD} = process.env
const DEPLOYMENTS_PATH = `${PWD}/deployments.json`
const CHAIN_CONFIG_PATH = `${PWD}/chain_config.json`

interface InnerObject {
  bridge: string;
  chains: Record<string, any>; 
  contracts: Record<string, any>; 
  links: { dst: string; src: string };
}

export class Deployments {
  map: Map<string, any>;
  private static instance: Deployments;

  constructor(map: any) {
    this.map = map;
  }

  public get(target: string) {
    return  this.map.get('chains')[target];
  }

  public getContracts(target: string) {
    return  this.map.get('contracts')[target];
  }

  public set(target: string, data: any) {
    this.map.set(target, data);
  }

  public getSrc() {
    // console.log(this.map.get('links'));
    return this.map.get('links').src;
  }

  public getDst() {
    return this.map.get('links').dst;
  }


  public static getDefault(config:string){

    if (!this.instance) {
      const data = fs.readFileSync(config);
      const json = JSON.parse(data.toString());
      const map = new Map(Object.entries(json));
       // Iterate over map entries to find the first entry with an object value
    let innerObject: InnerObject | undefined;;
    for (const [key, value] of map.entries()) {
      if (typeof value === 'object' && value !== null) {
        innerObject = value as InnerObject;
        break;
      }
    }

    if (innerObject) {
      // Create a new map with only the desired properties
      const filteredMap = new Map<string, any>([
        ['bridge', innerObject.bridge],
        ['chains', innerObject.chains],
        ['contracts', innerObject.contracts],
        ['links', innerObject.links]
      ]);

      this.instance = new this(filteredMap);
    }
    return this.instance;
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
      const data = fs.readFileSync("/home/dell/project/dive-core/cli/output/test1/dive_test1_61e3ad04453c.json");
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