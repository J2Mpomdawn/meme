const update = (() => {
  let current_sub = null;
  let guilds = {};
  let castles = [];
  const redraw = (castle_id) => {
    const info = castles[castle_id - 1];
    const node = document.querySelector(`gvg-castle[castle-id="${castle_id}"] > gvg-status`);
    const offset = [-3600*9, -3600*9, -3600*8, 3600*7, -3600, -3600][Math.floor(info.StreamId.WorldId / 1000) - 1];
    const fell = new Date(Math.max(0, info.UtcFallenTimeStamp+offset*1000));
    node.querySelector('gvg-status-icon-offense').innerText = info.AttackPartyCount;
    node.querySelector('gvg-status-icon-defense').innerText = info.DefensePartyCount;
    let defender = guilds[info.GuildId] || document.querySelector(`gvg-castle[castle-id="${castle_id}"] > gvg-castle-name`).innerText + _(' Forces');
    let attacker = guilds[info.AttackerGuildId] || '';
    if (info.GvgCastleState % 2 == 0) {
      node.removeAttribute('active');
      node.setAttribute('neutral', '');
    } else {
      node.removeAttribute('neutral');
      node.setAttribute('active', '');
    }
    if (info.GvgCastleState == 2 || info.GvgCastleState == 3) {
      [defender, attacker] = [attacker, defender];
    }
    node.querySelector('gvg-status-bar-offense').innerText = attacker;
    node.querySelector('gvg-status-bar-defense').innerText = defender;
  };
  const pack_le16 = (u8arr, offs, val) => {
    u8arr[offs] = val & 0xff;
    u8arr[offs + 1] = (val >>> 8) & 0xff;
  };
  const pack_le32 = (u8arr, offs, val) => {
    u8arr[offs] = val & 0xff;
    u8arr[offs + 1] = (val >>> 8) & 0xff;
    u8arr[offs + 2] = (val >>> 16) & 0xff;
    u8arr[offs + 3] = (val >>> 24) & 0xff;
  };
  const unpack_le16 = (u8arr, offs) => {
    return u8arr[offs] | u8arr[offs + 1] << 8;
  };
  const unpack_le32 = (u8arr, offs) => {
    return u8arr[offs] | u8arr[offs + 1] << 8 | u8arr[offs + 2] << 16 | u8arr[offs + 3] << 24;
  };
  const pack_stream_id = (u8arr, offs, w) => {
    const sid = w.WorldId << 19 | w.Class << 16 | w.GroupId << 8 | w.Block << 5 | w.CastleId;
    pack_le32(u8arr, offs, sid);
  };
  const unpack_stream_id = (u8arr, offset) => {
    const s = unpack_le32(u8arr, offset);
    return {
      value: {
        WorldId: s >>> 19,
        GroupId: (s >>> 8) & 0xFF,
        Class: (s >>> 16) & 0x7,
        Block: (s >>> 5) & 0x7,
        CastleId: s & 0x1F,
      },
      offset: offset + 4,
    };
  };
  const unpack_guild = (u8arr, offset, world) => {
    const gid = unpack_le32(u8arr, offset);
    const len = u8arr[offset + 4];
    const str = new TextDecoder('utf-8').decode(u8arr.subarray(offset + 5, offset + 5 + len));
    return {
      value: {
        GuildId: gid * 1000 + world % 1000,
        GuildName: str,
      },
      offset: offset + 5 + len,
    };
  };
  const unpack_castle = (u8arr, offset, world) => {
    const def_guild = unpack_le32(u8arr, offset);
    const atk_guild = unpack_le32(u8arr, offset + 4);
    const fallen = unpack_le32(u8arr, offset + 8) * 1000;
    const def_count = unpack_le16(u8arr, offset + 12);
    const atk_count = unpack_le16(u8arr, offset + 14);
    const state = u8arr[offset+16];
    return {
      value: {
        GuildId: unpack_le32(u8arr, offset) * 1000 + world % 1000,
        AttackerGuildId: unpack_le32(u8arr, offset + 4) * 1000 + world % 1000,
        UtcFallenTimeStamp: unpack_le32(u8arr, offset + 8) * 1000,
        DefensePartyCount: unpack_le16(u8arr, offset + 12),
        AttackPartyCount: unpack_le16(u8arr, offset + 14),
        GvgCastleState: u8arr[offset + 16],
      },
      offset: offset + 20,
    };
  };
  const check_stream_compat = (x, y) => {
    if ((x.GroupId == 0 && x.Class == 0 && x.Block == 0) ||
        (y.GroupId == 0 && y.Class == 0 && y.Block == 0)) {
      return x.WorldId == y.WorldId;
    } else {
      return x.GroupId == y.GroupId && x.Class == y.Class && x.Block == y.Block;
    }
  };
  let sub;
  const start = () => {
    sub = new WebSocket('wss://api.mentemori.icu/gvg');
    sub.binaryType = 'arraybuffer';
    sub.addEventListener('open', (event) => {
      document.getElementById('error').innerText = '';
      if (current_sub !== null) {
        let buffer = new ArrayBuffer(4);
        let u8arr = new Uint8Array(buffer);
        pack_stream_id(u8arr, 0, current_sub);
        sub.send(buffer);
      }
    });
    sub.addEventListener('message', (event) => {
      const u8arr = new Uint8Array(event.data);
      let offset = 0;
  
      while (offset < u8arr.length) {
        let res = unpack_stream_id(u8arr, offset);
        const stream_id = res.value;
        offset = res.offset;
  
        if (stream_id.CastleId == 0) {
          res = unpack_guild(u8arr, offset, stream_id.WorldId);
          const guild = res.value;
          guild.StreamId = stream_id;
          offset = res.offset;
  
          if (!check_stream_compat(stream_id, current_sub)) continue;
  
          if (guild.GuildId == 0) {
            guilds = {};
          } else {
            guilds[guild.GuildId] = guild.GuildName;
          }
        } else {
          res = unpack_castle(u8arr, offset, stream_id.WorldId);
          const castle = res.value;
          castle.StreamId = stream_id;
          offset = res.offset;
          if (!check_stream_compat(stream_id, current_sub)) continue;
  
          castles[stream_id.CastleId - 1] = castle;
          redraw(stream_id.CastleId);
        }
      }
    });
    sub.addEventListener('error', (event) => {
      console.log('Error', event);
      document.getElementById('error').innerText = 'WebSocket error';
    });
    sub.addEventListener('close', (event) => {
      document.getElementById('error').innerText = 'Connection closed, retrying in 5s';
      setTimeout(start, 5000);
    });
  };
  start();
  return (w) => {
    if (sub.readyState != WebSocket.OPEN) {
      current_sub = w;
      return;
    }
    if (current_sub) {
      // Unsubscribe
      let buffer = new ArrayBuffer(4);
      let u8arr = new Uint8Array(buffer);
      pack_stream_id(u8arr, 0, current_sub);
      sub.send(buffer);
    }
    current_sub = w;
    // Subscribe
    let buffer = new ArrayBuffer(4);
    let u8arr = new Uint8Array(buffer);
    pack_stream_id(u8arr, 0, current_sub);
    sub.send(buffer);
  };
})();