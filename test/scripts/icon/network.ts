import fs from 'fs';
import {IconService, Wallet} from 'icon-sdk-js';
import {ChainConfig} from "../setup/config";

const {IconWallet, HttpProvider} = IconService;

export class IconNetwork {
  iconService: IconService;
  nid: number;
  wallet: Wallet;
  private static instances: Map<string, IconNetwork> = new Map();

  constructor(_iconService: IconService, _nid: number, _wallet: Wallet) {
    this.iconService = _iconService;
    this.nid = _nid;
    this.wallet = _wallet;
  }

  public static getDefault() {
    return this.getNetwork('icon0');
  }

  public static getNetwork(target: string) {
    const entry = this.instances.get(target);
    if (entry) {
      return entry;
    }
    const config: any = ChainConfig.getChain(target);
    const httpProvider = new HttpProvider(config.endpoint);
    const iconService = new IconService(httpProvider);
    const keystore = this.readFile(config.keystore);
    const keypass = config.keysecret
      ? this.readFile(config.keysecret)
      : config.keypass;
    const wallet = IconWallet.loadKeystore(keystore, keypass, false);
    const nid = parseInt(config.network.split(".")[0], 16);
    const network = new this(iconService, nid, wallet);
    this.instances.set(target, network);
    return network;
  }

  private static readFile(path: string) {
    return fs.readFileSync(path).toString();
  }

  async getTotalSupply() {
    return this.iconService.getTotalSupply().execute();
  }

  async getLastBlock() {
    return this.iconService.getLastBlock().execute();
  }

  async getBTPNetworkInfo(nid: string) {
    return this.iconService.getBTPNetworkInfo(nid).execute();
  }

  async getBTPHeader(nid: string, height: string) {
    return this.iconService.getBTPHeader(nid, height).execute();
  }
}
