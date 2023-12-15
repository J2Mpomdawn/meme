package service

import (
	"fmt"
	"log"
	"strings"

	"meme/model"
)

func RegisterGuild(guilds map[int]model.Value_GuildId) {
	query := `
	insert into guilds (
		world_id,
		guild_id,
		guild_name
	) values %s
	on duplicate key update
		update_date = current_timestamp();
	`

	values := make([]string, 0, 20)
	for _, guild := range guilds {
		value := `
		(
			%d,
			%d,
			'%s'
		)
		`
		value = fmt.Sprintf(value,
			guild.StreamId.WorldId,
			guild.GuildId,
			guild.GuildName)
		values = append(values, value)
	}
	query = fmt.Sprintf(query, strings.Join(values, ","))
	query = strings.Join(strings.Fields(query), " ")

	exec_qery(query)
}

func RegisterRecord(castles map[int]model.Value_CastleId) {
	query := `
	insert into gvg_records (
		world_id,
		group_id,
		class,
		block,
		castle_id,
		def_guild_id,
		atk_guild_id,
		utc_fallen_time_stamp,
		def_count,
		atk_count,
		state
	) values %s
	on duplicate key update
		update_date = current_timestamp();
	`

	values := make([]string, 0, 20)
	for _, castle := range castles {
		if castle.StreamId.WorldId == 0 {
			continue
		}

		value := `
		(
			%d,
			%d,
			%d,
			%d,
			%d,
			%d,
			%d,
			%d,
			%d,
			%d,
			%d
		)
		`
		value = fmt.Sprintf(value,
			castle.StreamId.WorldId,
			castle.StreamId.GroupId,
			castle.StreamId.Class,
			castle.StreamId.Block,
			castle.StreamId.CastleId,
			castle.GuildId,
			castle.AttackerGuildId,
			castle.UtcFallenTimeStamp,
			castle.DefensePartyCount,
			castle.AttackPartyCount,
			castle.GvgCastleState)
		values = append(values, value)
	}
	query = fmt.Sprintf(query, strings.Join(values, ","))
	query = strings.Join(strings.Fields(query), " ")

	exec_qery(query)
}

func exec_qery(query string) {
	_, err := DbEngine.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
