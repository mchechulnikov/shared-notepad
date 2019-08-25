import {AfterViewInit, Component, HostListener, Input, OnInit} from '@angular/core';
import {Subject} from 'rxjs';
import {debounceTime, distinctUntilChanged} from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less']
})
export class AppComponent implements OnInit, AfterViewInit {
  private readonly socket: WebSocket;
  @Input() input = new Subject<string>();
  text = '';

  private static indent(withShift: boolean): void {
    if (withShift) {
      document.execCommand('outdent', true, null);
    } else {
      document.execCommand('indent', true, null);
    }
  }

  private static insert(str: string): void {
    document.execCommand('insertText', true, str);
  }

  constructor() {
    this.socket = new WebSocket('ws://localhost:5000/room/join');
  }

  ngOnInit(): void {
    this.socket.onmessage = (e) => {
      let value = e.data;
      value = value
        .split('\n')
        .map(x => '<div>' + (x ? x : '<br>') + '</div>').
        join('');
      this.text = value;
    };
  }

  ngAfterViewInit(): void {
    this.input
      .pipe(
        debounceTime(300),
        distinctUntilChanged()
      )
      .subscribe((value) => this.socket.send(value));
  }

  @HostListener('window:keydown', ['$event'])
  onKeyDown(event: KeyboardEvent) {
    const selection = document.getSelection().toString();
    switch (event.key) {
      // case 'Enter':
      //   //AppComponent.insert('\n');
      //   //document.execCommand('defaultParagraphSeparator', true, 'br');
      //   //document.execCommand('insertHTML',false,'<br>');
      //   //event.preventDefault();
      //   break;
      case 'Tab':
        AppComponent.indent(event.shiftKey);
        //AppComponent.insert('    ');
        event.preventDefault();
        break;
      case '(':
        AppComponent.insert(selection ? '(' + selection + ')' : '()');
        event.preventDefault();
        break;
      case '[':
        AppComponent.insert(selection ? '[' + selection + ']' : '[]');
        event.preventDefault();
        break;
      case '{':
        AppComponent.insert(selection ? '{' + selection + '}' : '{}');
        event.preventDefault();
        break;
    }
  }

  onChange(e) {
    let value = e.target.innerHTML;

    // Convert `&amp;` to `&`.
    value = value.replace(/&amp;/gi, '&');

    // Replace spaces.
    value = value.replace(/&nbsp;/gi, ' ');

    // Remove "<b>".
    value = value.replace(/<b>/gi, '');
    value = value.replace(/<\/b>/gi, '');

    // Remove "<strong>".
    value = value.replace(/<strong>/gi, '');
    value = value.replace(/<\/strong>/gi, '');

    // Remove "<i>".
    value = value.replace(/<i>/gi, '');
    value = value.replace(/<\/i>/gi, '');

    // Remove "<em>".
    value = value.replace(/<em>/gi, '');
    value = value.replace(/<\/em>/gi, '');

    // Remove "<u>".
    value = value.replace(/<u>/gi, '');
    value = value.replace(/<\/u>/gi, '');

    // Replace "<div>" (from Chrome).
    value = value.replace(/<div>/gi, '\n');
    value = value.replace(/<\/div>/gi, '');
    value = value.replace(/^\n/gi, '');

    // Replace "<p>" (from IE).
    value = value.replace(/<p>/gi, '\n');
    value = value.replace(/<\/p>/gi, '');

    this.input.next(value);
  }
}
