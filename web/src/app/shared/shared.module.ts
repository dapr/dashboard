import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MonacoEditorModule } from 'ng-monaco-editor';
import { EditorComponent } from './editor/editor.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    MonacoEditorModule
  ],
  declarations: [
    EditorComponent
  ],
  exports: [
    EditorComponent
  ]
})
export class SharedModule { }
