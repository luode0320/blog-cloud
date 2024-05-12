interface RsaCondition {
  message?: string;
  bits?: number;
  isPKCS8?: boolean;
  privateKey?: string;
  publicKey?: string;
  sign?: string;
  signType?: string;
}

interface RsaResult {
  privateKey: string;
  publicKey: string;
}
