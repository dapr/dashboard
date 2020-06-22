import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { DaprComponentsComponent } from './dapr-components.component';
import { MatExpansionModule } from '@angular/material/expansion';
import { CommonModule } from '@angular/common';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    MatExpansionModule,
  ],
  declarations: [
    DaprComponentsComponent,
  ],
})
export class DaprComponentsModule { }
