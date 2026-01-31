export const ROOT_INDEX = "/";
export const ROOT_ADMIN = "/admin";

export const PATH_ADMIN_LOGIN = "login";
export const PATH_ADMIN_DASHBOARD = "dashboard";

export const getPathWithRoot = (root: string, path: string) =>
  `${root}/${path}`;
