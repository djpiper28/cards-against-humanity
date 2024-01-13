import { expect, test } from "vitest";
import { RpcMessage } from "./rpcTypes";

test("RPC types should be defined", () => {
  const testMessage: RpcMessage = {};
  expect(testMessage).toBeDefined();
});
