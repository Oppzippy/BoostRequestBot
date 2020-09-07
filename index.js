const Discord = require('discord.js');
const client = new Discord.Client();
const token = 'NzQ5MDM5MTIzODkzMDU5NjI1.X0mLCw.DvEtdNs-jHaf8czmdfkZizAdQJc';
const BOOST_REQUEST_BACKEND_CHANNEL_ID = '751971313668718705';
const BOOST_REQUEST_CHANNEL_ID = '719639847316357180';
const reactionArray = ['👎', '👍', '💩', '❌'];
const boostRequestsBySignupMessageId = new Map();

client.login(token);
client.on('ready', () => {
	console.log('client online!');

});

// TODO: Edit first embed to get of 2nd embed
// TODO: Clean up embeds
// TODO: Fix bug so it doesnt respond to one person and mark all with checks
// TODO: Put every reaction.message.id into an array and check again the current one being operated on so we know if we should be using it or not


// Event Catcher when users react to messages
client.on('messageReactionAdd', async (reaction, user) => {
	// if user is not the bot + reaction was in backend channel + Confirm reaction is Thumbs Up
	const signupMessage = boostRequestsBySignupMessageId.get(reaction.message.id);
	const guildMember = reaction.message.guild.members.cache.get(user.id);

	if (signupMessage && !user.bot && reaction.emoji.name === '👍' && !reaction.message.reactions.cache.has('✅')) {
		// TODO: Allow all advertisers to react after two minutes.
		// Don't forget to look over existing reactions again in case
		// a regular advertiser reacted before the two minutes.
		if (guildMember.roles.cache.some(role => role.name === 'Elite Advertiser')) {
			boostRequestsBySignupMessageId.delete(reaction.message.id);
			try {
				const winnerName = user.username;
				console.log(winnerName + ' won!');
				// remove reactions.
				await reaction.message.reactions.removeAll().catch(error => console.error('Failed to clear reactions: ', error));
				await reaction.message.react('✅');
				console.log(reaction.message.reactions.cache.has('✅'));
				sendEmbed(user, signupMessage.requesterId, winnerName);
			}
			catch (error) {
				console.error('One of the emojis failed to react.');
			}
		}
	}
});

// Event Catcher when users send a message
client.on('message', async message => {
	console.log(message.content);
	// If User is not a bot AND is messsaging in BoostRequest Channel
	if (!message.author.bot && message.channel.id === BOOST_REQUEST_CHANNEL_ID) {
		// Create embed in the Backend Channel
		const signupMessage = await BREmbed(message);
		boostRequestsBySignupMessageId.set(signupMessage.id, {
			requesterId: message.author.id,
			createdAt: message.createdAt,
		});
	}
});


async function BREmbed(brMessage) {
	// Variable to eaily add hyperlink to the original message.
	const ref = 'https://discordapp.com/channels/' + brMessage.guild.id + '/' + brMessage.channel.id + '/' + brMessage.id;
	const messagelink = '[' + brMessage.content + '](' + ref + ')';
	const exampleEmbed = new Discord.MessageEmbed()
		.setColor('#0000FF')
		.setTitle('New Boost Request')
		.setThumbnail(brMessage.author.displayAvatarURL())
		.addFields({ name: brMessage.author.username, value: messagelink })
		.setTimestamp()
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
	// Send embed to BoostRequest backend THEN add the Thumbsup Icon
	const message = await client.channels.cache.get(BOOST_REQUEST_BACKEND_CHANNEL_ID).send(exampleEmbed);
	shuffle(reactionArray);
	const reactPromises = reactionArray.map(emoji => message.react(emoji));
	await Promise.all(reactPromises);
	return message;
}


function sendEmbed(embedUser, requesterId, winnerName) {
	// Make Embed post here
	console.log(winnerName + 'won!');
	const selectionBRBEmbed = new Discord.MessageEmbed()

		.setColor('#FF0000')
		.setTitle('And the fastest clicker is...')
		.setThumbnail(embedUser.displayAvatarURL())
		.addFields(
			{ name: embedUser.username, value: 'Please message <@' + requesterId + '>' })
		.setTimestamp()
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
	client.channels.cache.get(BOOST_REQUEST_BACKEND_CHANNEL_ID).send(selectionBRBEmbed);

	// Make Embed post here
	const selectionBREmbed = new Discord.MessageEmbed()
		.setColor('#00FF00')
		.setTitle('Huokan Boosting Community Boost Request')
		.setThumbnail('https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png')
		.addFields(
			{ name: 'Your advertiser has been chosen.', value:'They will message you shortly <@' + requesterId + '>.' })
		.setTimestamp()
		.setFooter('Huokan Boosting Community', 'https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png');
	client.channels.cache.get(BOOST_REQUEST_CHANNEL_ID).send(selectionBREmbed);
}

function shuffle(array) {
	array.sort(() => Math.random() - 0.5);
}