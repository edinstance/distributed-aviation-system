export type AuthContext = {
  organization: {
    id: string;
    name: string;
    schemaName: string;
  };
  admin: {
    id: string;
    username: string;
    password: string;
    email: string;
  };
};
