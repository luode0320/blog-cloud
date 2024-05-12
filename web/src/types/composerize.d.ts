declare module "composerize" {
  const composerize: (dockerRun: string, existCompose: string) => string;
  export default composerize;
}
