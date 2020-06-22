import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConfigurationComponent } from './configuration.component';
import { MatTableModule } from '@angular/material/table';

@NgModule({
  declarations: [ConfigurationComponent],
  imports: [
    CommonModule,
    MatTableModule
  ]
})
export class ConfigurationModule { }
