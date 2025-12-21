import type { GetAuthSchema } from "../../../gen/api/v1/auth_pb"

export const methods: Record<keyof typeof GetAuthSchema.field, string> = {
  "authless": "Authless",
  "basicAuth": "Basic Auth",
  'proxyAuth': "Proxy Auth"
}

export const getMethod = (key: keyof typeof methods | undefined): string => {
  if (key) {
    return methods[key] ?? "Unknown method"
  }
  return "Unknown method"

}
