import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';

import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';

import { NotificationsService } from 'angular2-notifications';

import { APIResponse } from '../common/api-response';

import { Member } from "./member";
import { MemberService } from './member.service';

@Component({
    templateUrl: 'app/members/member.component.html',
})
export class MemberComponent {
    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private memberService: MemberService,
        private notificationService: NotificationsService
    ) { }

    member: Member = new Member();
    errors: string[] = [];

    ngOnInit(): void {
        this.route.params.forEach((params: Params) => {
            if (params['id'] == 'new') {
                this.member = new Member;
                return;
            }

            let id = +params['id'];
            if (id != this.member.id) {
                this.memberService.getMember(id)
                    .then((member: Member) => this.member = member)
                    .catch((error: any) => alert(error));
            }
        });
    }

    private save() {
        let call: Promise<APIResponse>;
        if (this.member.id == undefined) {
            call = this.memberService.postMember(this.member);
        } else {
            call = this.memberService.putMember(this.member.id, this.member);
        }

        call.then(response => {
            if (response.errors.length > 0) {
                response.errors.forEach((error: any) => {
                    switch (error) {
                        default:
                            alert('unhandled error: ' + error);
                            break;
                    }
                });
                this.errors = response.errors;
                return;
            }
            this.errors = [];
            if (this.member.id == undefined) {
                this.member.id = response.id;
                this.router.navigate(['/stock', 'member', response.id]);
                this.notificationService.create('New member created', `${this.member.initials} ${this.member.lastName} has been created.`, 'success');
            } else {
                this.notificationService.create('Changes saved', `Changes to ${this.member.initials} ${this.member.lastName} were saved.`, 'success');
            }
        })
            .catch(error => alert(error));
    }
}
