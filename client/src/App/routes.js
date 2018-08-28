import {OrdersPage, DealsPage, HomePage, LoginPage, SummaryPage} from '../_pages';

export const publicRoutes = [{
  component: HomePage,
  path: '/'
}, {
  component: LoginPage,
  path: '/login'
}];

export const privateRoutes = [{
  component: DealsPage,
  path: '/deals'
}, {
  component: OrdersPage,
  path: '/bids'
}, {
  component: SummaryPage,
  path: '/summary'
}];