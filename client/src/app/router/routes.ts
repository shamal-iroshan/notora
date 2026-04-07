export const ROOT_INDEX = "/";
export const ROOT_ADMIN = "/admin";

export const PATH_ADMIN_LOGIN = "login";
export const PATH_ADMIN_DASHBOARD = "dashboard";

export const PATH_CLIENT_LOGIN = "login";
export const PATH_CLIENT_SIGN_UP = 'sign-up';
export const PATH_CLIENT_SIGN_UP_SUCCESS = 'sign-up-success'
export const PATH_CLIENT_NOTES = "notes";
export const PATH_CLIENT_ACCOUNT = "account";

export const getPathWithRoot = (root: string, path: string) =>
  `${root}/${path}`;
