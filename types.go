package main

import ()

type Member struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Active            bool   `json:"active"`
	Subscription_type string `json:"subscription_type"`
	Join_date         string `json:"join_date"`
}

func newMember(member Member) *Member {
	return &Member{
		Name:              member.Name,
		Active:            member.Active,
		Subscription_type: member.Subscription_type,
		Join_date:         member.Join_date,
	}
}
