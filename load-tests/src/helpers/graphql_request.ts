import http, { RefinedResponse, ResponseType } from "k6/http";
import { check, JSONArray } from "k6";
import { TypedDocumentNode } from "@graphql-typed-document-node/core";
import { print } from "graphql";

export function graphql<
  TData = any,
  TVariables extends Record<string, unknown> = Record<string, unknown>,
>(
  url: string,
  doc: TypedDocumentNode<TData, TVariables>,
  variables?: TVariables,
): TData {
  const res: RefinedResponse<ResponseType> = http.post(
    url,
    JSON.stringify({
      query: print(doc),
      variables,
    }),
    { headers: { "Content-Type": "application/json" } },
  );

  check(res, {
    "status is 200": (r) => r.status === 200,
    "response has JSON": (r) => {
      try {
        r.json();
        return true;
      } catch {
        return false;
      }
    },
    "no GraphQL errors": (r) => {
      const errors = r.json("errors") as JSONArray | null;
      return !errors || errors.length === 0;
    },
  });

  if (res.json("errors")) {
    console.error("GraphQL errors:", res.json("errors"));
    console.error("Response body:", res.body);
  }

  return res.json("data") as TData;
}
