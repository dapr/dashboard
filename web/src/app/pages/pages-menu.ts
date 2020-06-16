export const MENU_ITEMS: MenuItem[] = [
  {
    title: 'Dashboard',
    icon: 'home-outline',
    link: '/dashboard',
    home: true,
  },
];

export const COMPONENTS_MENU_ITEM = {
  title: 'Components',
  icon: 'keypad-outline',
  link: '/components',
  home: false,
  name: 'components',
};

export interface MenuItem {
  title: string,
  icon: string,
  link: string,
  home: boolean
}
