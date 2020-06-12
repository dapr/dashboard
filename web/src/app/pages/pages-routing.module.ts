import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';

import { PagesComponent } from './pages.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ComponentsComponent } from './components/components.component';
import { LogsComponent } from './logs/logs.component';
import { DetailComponent } from './detail/detail.component';

const routes: Routes = [{
  path: '',
  component: PagesComponent,
  children: [
    {
      path: 'dashboard',
      component: DashboardComponent,
    },
    {
      path: 'components',
      component: ComponentsComponent,
    },
    {
      path: 'logs/:id',
      component: LogsComponent,
    },
    {
      path: 'detail/:id', 
      component: DetailComponent,
    },
    {
      path: '',
      redirectTo: 'dashboard',
      pathMatch: 'full',
    },
    { path: '**', redirectTo: 'dashboard' },
  ],
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PagesRoutingModule { }