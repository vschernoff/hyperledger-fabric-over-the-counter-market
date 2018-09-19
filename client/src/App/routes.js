import {HomePage, LoginPage, SummaryPage} from '../_pages';

export const publicRoutes = [{
  component: HomePage,
  path: '/'
}, {
  component: LoginPage,
  path: '/login'
}];

export const privateRoutes = [{
  component: SummaryPage,
  path: '/summary'
}];
