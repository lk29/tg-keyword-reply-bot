<h1 align="center">
   <br>Telegram keyword auto reply robot<br>
</h1>

### Basic commands

- Add keyword reply rule `/add keyword===reply content` or `/add keyword1||keyword2===reply content`
- Keywords can use regular expressions, such as `/add re:p([a-z]+)ch===test regular`, it will match the rule `p([a-z]+)ch`
- The keyword deletion rule `/del keyword` does not support deleting multiple keywords at once
- Automatically delete text messages containing keywords, just set the reply content to `delete`, and add delete message permissions to the robot
- Use `/list` command to view all auto-reply rules in this group
- Add the management authority to delete messages and kick people to the robot, which can automatically prevent Halal (Arabic)

### Reply with special content

- Reply content supports text\picture\GIF\video, default text
- If you need a picture, set the reply content to `photo: https://t.me/c/1472018167/53095`, `https://t.me/c/1472018167/53095` is the picture that has been sent to get it the link to
- Similarly, for gif, replace `photo` with `gif`, replace video with `video`, and replace files with `file`
- Note: The link here must be a public group, otherwise it cannot be sent out

## License

MIT zu1k i@zu1k.com

Do not provide any robot-related consulting or other services, do not chat with me privately
