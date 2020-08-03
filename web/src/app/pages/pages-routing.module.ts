import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { PagesComponent } from './pages.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { DaprComponentsComponent } from './dapr-components/dapr-components.component';
import { DetailComponent } from './dashboard/detail/detail.component';
import { DaprComponentDetailComponent } from './dapr-components/dapr-component-detail/dapr-component-detail.component';
import { ConfigurationComponent } from './configuration/configuration.component';
import { ConfigurationDetailComponent } from './configuration/configuration-detail/configuration-detail.component';
import { ControlPlaneComponent } from './controlplane/controlplane.component';

const routes: Routes = [{
  path: '',
  component: PagesComponent,
  children: [
    {
      path: 'overview',
      component: DashboardComponent,
    },
    {
      path: 'overview/:id',
      component: DetailComponent,
    },
    {
      path: 'components',
      component: DaprComponentsComponent,
    },
    {
      path: 'components/:name',
      component: DaprComponentDetailComponent,
    },
    {
      path: 'configurations',
      component: ConfigurationComponent,
    },
    {
      path: 'configurations/:name',
      component: ConfigurationDetailComponent,
    },
    {
      path: 'controlplane',
      component: ControlPlaneComponent,
    },
    {
      path: '',
      redirectTo: 'overview',
      pathMatch: 'full',
    },
    { path: '**', redirectTo: 'overview' },
  ],
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PagesRoutingModule { }
