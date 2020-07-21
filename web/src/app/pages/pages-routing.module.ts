import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { PagesComponent } from './pages.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { DaprComponentsComponent } from './dapr-components/dapr-components.component';
import { DetailComponent } from './dashboard/detail/detail.component';
import { DaprComponentDetailComponent } from './dapr-components/dapr-component-detail/dapr-component-detail.component';
import { ConfigurationComponent } from './configuration/configuration.component';
import { ControlPlaneComponent } from './controlplane/controlplane.component';

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
      component: DaprComponentsComponent,
    },
    {
      path: 'components/:id',
      component: DaprComponentDetailComponent,
    },
    {
      path: 'detail/:id',
      component: DetailComponent,
    },
    {
      path: 'configuration',
      component: ConfigurationComponent,
    },
    {
      path: 'controlplane',
      component: ControlPlaneComponent,
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
