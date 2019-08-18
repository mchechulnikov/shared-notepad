import {Component, HostListener, ViewChild} from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less']
})
export class AppComponent {
  // @ts-ignore
  @ViewChild('editor') editorEl;

  @HostListener('window:keydown', ['$event'])
  onKeyDown(event: KeyboardEvent) {
    const selection = document.getSelection().toString();
    switch (event.key) {
      case 'Tab':
        AppComponent.Indent(event.shiftKey);
        event.preventDefault();
        break;
      case '(':
        AppComponent.Insert(selection ? '(' + selection + ')' : '()');
        event.preventDefault();
        // debugger;
        this.editorEl.nativeElement.focus();
        this.editorEl.nativeElement.setSelectionRange(1, 1);
        break;
      case '[':
        AppComponent.Insert(selection ? '[' + selection + ']' : '[]');
        event.preventDefault();
        break;
      case '{':
        AppComponent.Insert(selection ? '{' + selection + '}' : '{}');
        event.preventDefault();
        break;
    }
  }

  private static Indent(withShift: boolean): void {
    if (withShift)
      document.execCommand('outdent', true, null);
    else
      document.execCommand('indent', true, null);
  }

  private static Insert(str: string): void {
    document.execCommand('insertText', true, str);
  }
}
