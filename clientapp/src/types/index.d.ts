export type MySessionData = {
  email: string;
  fullName: string;
  hasOwnerRole: boolean;
  tenantName: string;
  customTenantDomain: string;
  now: string;
}

export type ApiSessionData = {
  email: string;
  fullName: string;
  tenantName: string;
  customTenantDomain: string;
  now: string;
  roles: {
    id: string;
    name: string;
    displayName: string;
  }[];
}
