import {AfterViewInit, Component, HostListener, Input, OnInit, ViewChild} from '@angular/core';
import {Observable, Subject} from 'rxjs';
import {debounceTime, distinctUntilChanged, subscribeOn} from "rxjs/operators";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less']
})
export class AppComponent implements OnInit, AfterViewInit {
  private readonly socket: WebSocket;
  @Input() input = new Subject<string>();
  text = '';

  // @ts-ignore
  @ViewChild('editor') editorEl;

  // private static indent(withShift: boolean): void {
  //   if (withShift) {
  //     document.execCommand('outdent', true, null);
  //   } else {
  //     document.execCommand('indent', true, null);
  //   }
  // }

  private static insert(str: string): void {
    document.execCommand('insertText', true, str);
  }

  constructor() {
    this.socket = new WebSocket('ws://localhost:5000/room/join');
  }

  ngOnInit(): void {
    this.socket.onmessage = (e) => this.text = e.data;
  }

  ngAfterViewInit(): void {
    this.input
      .pipe(
        debounceTime(500),
        distinctUntilChanged()
      )
      .subscribe((value) => this.socket.send(value));
  }

  @HostListener('window:keydown', ['$event'])
  onKeyDown(event: KeyboardEvent) {
    const selection = document.getSelection().toString();
    switch (event.key) {
      case 'Tab':
        AppComponent.insert('    ');
        event.preventDefault();
        break;
      case '(':
        AppComponent.insert(selection ? '(' + selection + ')' : '()');
        event.preventDefault();
        this.editorEl.nativeElement.focus();
        this.editorEl.nativeElement.setSelectionRange(1, 1);
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
}
