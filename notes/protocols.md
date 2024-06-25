# Protocols

Coming from a background in web development, I've always taken HTTP for granted. I could regurgitate facts about it: it's a protocol, a set of rules—easy peasy. And I could work with it just fine. That's the value of black boxing some complex topic—you can drive a car without knowing the specifics of internal combustion.

But now that I'm beginning to implement various protocols myself, I'm forced to confront the fact that I've never really given them much consideration beyond being aware of their existence.

I wanted to write down some bits of intuition I've hobbled together while journeying down this rabbit hole.

## Example: Two Prisoners

Let's say you have two prisoners who can neither see nor hear one another. The only thing connecting them is a bit of spare rope sprawled across the cement floor haphazardly.

If they wanted to communicate, they’d have to use this rope. Great. How are they going to do that?

### Initial Attempt

Suppose they’d agreed beforehand to swing the rope up and down to cause a wave. Then, using the alphabet to their advantage, they’d use a simple encoding where the letter 'a' was 1, and the letter 'z' was 26. Fine.

Well... not fine. They’re going to quickly run into a load of questions:

- How do you know when to stop counting the number of waves for each letter? For example, if you receive 3 waves, that could be: 'ab' or 'c'.
- How do you know when the message has stopped and there isn’t just a delay?
- What if the other person wants to reply? If they reply using the same rope, there will be collisions.
- What if the wind blows the rope a bit causing _random noise_, could you filter that out somehow?

### Need for Structure

It seems they need a more structured way to communicate here.

_Enter the protocol_

Defining the Basic Unit of Information
These problems generalize to any form of communication and any medium in which that communication occurs. One of the first things you may need to decide on is: what is the basic unit of information?

We said that these two prisoners used a single wave to represent a unit of communication. With this, you’d have to count each wave and assign it to some value.

It works, but that’s just SO many waves to count. The prisoners have to count all those waves. What if they lose count? The more waves you have to send to represent a character, the more opportunities for errors there are.

Wouldn't it be nicer if we didn’t have to send so many just to represent a single character?

### Introducing Multiple States

The problem is that our unit of communication can only represent one state. What if instead of a single wave of the rope, we just agreed that we could send one of two different types of waves: a _really big_ one and a _really small_ one? The emphasis on the word _really_ is important. If the difference between the two isn’t large enough, we may start guessing if we sent a small wave or a large one, thereby introducing _ambiguity_. We can't have that because now the prisoners would be stuck interpreting the amplitude of the waves.

Okay, so let's say that a wave with a height of a foot or more from bottom to top is big and one that is between 1 to 2 inches is small.

What does this afford us?

### Advantages Of Multiple States

With this new idea, we can tremendously cut down on the number of waves we'd need to send for our message. Let's still map the letters 'a' through 'z' in the same way we did before, with 'z' being 26. But let's say that we'll map a _big_ wave to the number 1 and the _small_ wave to the number 0. We... did we just reinvent binary? Yeah, I suppose we did.

One benefit of this is increased efficiency. Instead of sending 26 waves to represent the letter 'z', we can send this wave sequence:

```
big wave, big wave, small wave, big wave, small wave
```

Or in other terms:

```
11010
```

Well, in binary, this represents the number _26_ which is... the letter 'z'. So wait a minute. We only have to send 5 waves to represent something that used to take us 26 waves? Yep. That's the power of being able to represent multiple states.

That's awesome!

Now you may be tempted to say something like: well what about 3 states? Trinary? That would cut down the amount of waves we'd have to send for sure! And you'd be right. If we had 3 states: Big wave(2), medium wave (1) and small wave (0), we'd only have to send THREE waves to represent the number 26!

Example:

```
Big wave (2), Big wave (2), Big wave (2)
```

This is basically saying you have two 9s (18), two 3's (6), and two 1s (2). The end result would be adding these values together. 18 + 6 + 2 = 26! Isn't math _fun_?

But we'd ultimately introduce more ambiguity, because with trinary and the addition of this new _medium wave_, the prisoners may have to really be paying attention to the size of the waves coming across. This increased need for measurement can introduce bugs and complicate things. So we'll stick with binary for now.

### Data About Data

This is really great! The prisoners can now send less waves to communicate with one another. But there are still questions lingering.

- How does the other prisoner know when the message has stopped?
- How big is the message? Will the prisoner need to write it down because it's so large? This is useful information.
- Is the message still English or are there other things that can be communicated, like visuals? Pictures? Music?

I mean, look they probably just want to chit chat. They're not sending PDFs over a bit of spare rope. But let's stick with this example. Would it be possible to send information like that? Absolutely. Now whether or not the time would be reasonable for them is another story, but let's not worry about that for now.

But let's say the prisoners came up with a way to send a visual.

I don't want to go into too much detail here because quite frankly, it's over my head but this simple example should suffice.

Suppose the prisoners each had a piece of paper with a grid of dots that they could reuse. To draw a picture on the grid, you can mark points and then connect the dots. Fine. So the way that the prisoners can send an image is to specify the rows and columns for the dots to mark, and at the end, one only needs to connect them all to see the image. That should work.

But we're immediately going to have a problem.

_How do we know that the other prisoner is sending a picture? What if they're sending a message in English?_

You see, a stream of waves is inherently _ambiguous_ without a _way to assign meaning_.

So what should they do? Well suppose that they revised their protocol to handle this? With the protocol as is, the maximum character you can send is 'z' and that requires only five waves to represent. What if they agreed to set aside the first five waves in a message as a special _message about the message_?

So if they want to send a message in English, perhaps they send an 'e' using the first five waves, and a 'p' if they intend to send a picture?

Now you don't need the entire five waves to send the letter 'e' because we've encoded that as the number 5. Which can be expressed as `101` in binary, or `Big wave, Small wave, Big wave` in our particular protocol. So there will be some wasted waves, but let's forget about that for now.

Great, so remembering that the letter 'p' is encoded as 16, to send a picture message, you may get something like this in the beginning:

Decimal: 16
Binary: `10000`
Our protocol: `Big wave, Small wave, Small wave, Small wave, Small wave`

Now we can understand that someone is sending a picture. So any waves _after_ this must _contribute to the picture_.

Fantastic! We've just sent a picture across a bit of rope!

This information that came in the beginning was special. It was _data about the data_. We refer to this as a _header_. It's data that comes _ahead_ of the actual data.

Headers are important because again, a stream of waves or _bits_ is inherently ambiguous as we've said. Headers can _disambiguate_ and provide clarity.

To risk another analogy, whenever you send a letter (do people still do this?) you put it in an envelope right? And that envelope must contain data like: where to send the letter, who it's addressed to, who it's coming from and most important - postage!

Well it's no different here! Our clever prisoners will undoubtedly have invented something like this.

Another piece of data that's important for our prisoners to know would be: _how long is the incoming message_? I mean, are they gonna have to get a pen and paper out or is it short enough to where they could just figure it out?

Well how would solve this? Again, you have to think about the word: _convention_. These rules are agreed upon at some point either before hand or using the shoddy communication they have through an evolutionary process of getting pissed off because of errors in transmission and translation.

Okay, so suppose the prisoners said, what if, after the header for the _type_ of message, we allocate a certain number of waves for the _length_ of the message?

That sounds reasonable!

Let's say that you allocate 8 waves to represent the length of the incoming message. Again, this is binary we're talking here. So we need to see what is the maximum number that we could represent using 8 bits? Well if you have a 1 in every slot you'd have this number: `11111111`. That number in decimal is 255. Or you can remember the equation (2^N-1) where N is the number of bits. Why subtract 1? Yeah, that's because you have to also represent 0. If you didn't represent 0 you could represent the number 256.

Alright, so if we allocate 8 bits - a _byte_ of data (8 waves in our example), the maximum length of the message can be 255 waves long. If each character uses up multiple waves, this doesn't leave us with very much room for our emssage.

So we have two options, we could allocate _another byte_ or another 8 waves for the length of the message, or we could agree that messages over a certain length need to be split up.

Let's go with the latter for now. So breaking up a message eh? What does that look like? Well, if you had to send a series of messages to someone through the mail, you may segment the message and have something in the corner that says something like: 1 of 10. Where you have ten total letters and this is just the first!

That seems to work, and also, now you have an ordering so that if the letters arrive _out of order_, you can rearrange them at the end! But for the purposes of keeping this simple, let's be content with NOT implementing this ordering of letters for now.

So a length without a unit is again (sounding like a broken record here) _ambiguous_, so let's say that this will be the length in bytes or collections of 8 waves. Yes, there will be wasted waves, but that's tomorrow's worry.

Okay, that sounds good, so now a possible English message looks like this over the rope:

_I'm including a pipe character to make it easier to read, it's not part of the message_

**Protocol**

```
message type | length of message | message
```

**Decimal**:

```
5 | 2 | 8 9
```

**Binary**:

```
00101 | 00000010 | 00001000 00001001

```

**Interpretation**:

```
e (for English message) | 2 bytes long (or 16 waves) | hi
```

Hey now. That looks promising! Well you can see how useful these protocol things can be.

### Protocols Can Be Nested

Alright, so the headers have been taken care of. And that's wonderful. But we have these message types that may have _differing needs_. Pictures may have different headers than an English message. Messages in English may need to be treated differently.

Well, it seems like the protocol that deals with how we send any information at all isn't dependent upon these other more specific protocols like what headers to send right?

That's how it is in the real world too. Different protocols serve different purposes but some protocols _require others_ or we say they're built _on top of other protocols_.

An http request is a protocol that sits on top of a lower level protocol called _tcp_. TCP is what's known as a Transport Protocol. It deals with how to send messages between networks and specifically TCP deals with maintaining connections.
