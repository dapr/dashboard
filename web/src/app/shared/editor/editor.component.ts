import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { MonacoEditorOptions } from 'ng-monaco-editor';
import { ThemeService } from '../../theme/theme.service';

@Component({
  selector: 'app-editor',
  templateUrl: 'editor.component.html',
  styleUrls: ['editor.component.scss']
})
export class EditorComponent implements OnInit {

  @Input() model = '';
  @Input() language = '';
  @Output() modelChange = new EventEmitter<string>();

  options: MonacoEditorOptions = {
    folding: true,
    minimap: { enabled: true },
    readOnly: true,
    language: this.language || 'yaml',
    contextmenu: false,
    scrollBeyondLastLine: false,
    lineNumbers: false as any,
    theme: this.isDarkTheme() ? 'vs-dark' : 'vs',
    automaticLayout: true
  };

  constructor(
    private themeService: ThemeService,
  ) { }

  ngOnInit() {
    this.themeService.themeChanged.subscribe((newTheme: string) => {
      this.options = {
        ...this.options,
        theme: newTheme.includes('dark') ? 'vs-dark' : 'vs'
      };
    });
  }

  isDarkTheme(): boolean {
    return this.themeService.getTheme().includes('dark');
  }
}
