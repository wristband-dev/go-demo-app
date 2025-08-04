export type MySessionData = {
  email: string;
  fullName: string;
  hasOwnerRole: boolean;
  tenantDomainName: string;
}

export type ApiSessionData = {
  email: string;
  fullName: string;
  tenantDomainName: string;
  roles: {
    id: string;
    name: string;
    displayName: string;
  }[];
}
