export interface MenuItem {
  title: string,
  icon: string,
  link: string,
  name: string,
}

export const DASHBOARD_MENU_ITEM = {
  title: 'Overview',
  icon: 'home',
  link: '/dashboard',
  name: 'dashboard',
};

export const COMPONENTS_MENU_ITEM = {
  title: 'Components',
  icon: 'apps',
  link: '/components',
  name: 'components',
};

export const CONFIGURATIONS_MENU_ITEM = {
  title: 'Configurations',
  icon: 'build',
  link: '/configuration',
  name: 'configurations',
};

export const CONTROLPLANE_MENU_ITEM = {
  title: 'Control Plane',
  icon: 'settings',
  link: '/controlplane',
  name: 'status',
};

export var MENU_ITEMS: MenuItem[] = [
  DASHBOARD_MENU_ITEM,
];
