import { UserPermission } from '../gen/api/v1/auth_pb';
import type { Authenticated } from '../gen/api/v1/auth_pb';

export type AuthResult = Authenticated | false | undefined;

export function hasPermission(auth: AuthResult, ...allowed: UserPermission[]) {
  if (!auth) return false;
  return allowed.includes(auth.permission);
}

export function isCreatator(auth: AuthResult, owner?: string) {
  if (!auth) return false
  return auth.userId === owner
}
