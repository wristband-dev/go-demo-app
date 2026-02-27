export type MySessionData = {
  email: string;
  fullName: string;
  hasOwnerRole: boolean;
  tenantName: string;
  tenantCustomDomain: string;
  now: string;
}

export type ApiSessionData = {
  email: string;
  fullName: string;
  tenantName: string;
  tenantCustomDomain: string;
  now: string;
  roles: {
    id: string;
    name: string;
    displayName: string;
  }[];
}
